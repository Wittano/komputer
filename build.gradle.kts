import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.10"
    kotlin("kapt") version "1.9.0"

    application
}

group = "com.wittano"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

val picocliVersion = "4.7.4"

dependencies {
    implementation("com.discord4j:discord4j-core:3.2.5")
    implementation("org.apache.logging.log4j:log4j-core:2.20.0")
    implementation("ch.qos.logback:logback-classic:1.4.11")
    implementation("info.picocli:picocli:$picocliVersion")

    kapt("info.picocli:picocli-codegen:$picocliVersion")

    testImplementation(kotlin("test"))
}

kapt {
    arguments {
        arg("project", "${project.group}/${project.name}")
    }
}

tasks.test {
    useJUnitPlatform()
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}

application {
    mainClass.set("MainKt")
}