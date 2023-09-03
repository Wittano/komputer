import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.10"
    kotlin("kapt") version "1.9.10"
}

group = "com.wittano"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

val picocliVersion = "4.7.4"

dependencies {
    // Internal dependencies
    implementation(project(":core"))

    // Picocli
    implementation("info.picocli:picocli:$picocliVersion")
    kapt("info.picocli:picocli-codegen:$picocliVersion")

    // Logger
    implementation("ch.qos.logback:logback-classic:1.4.11")
    implementation("org.codehaus.janino:janino:3.1.10")

    testImplementation(platform("org.junit:junit-bom:5.9.1"))
    testImplementation("org.junit.jupiter:junit-jupiter")
}

// TODO Create optional native image with GrallVM

tasks.test {
    useJUnitPlatform()
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}