import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.10"
    kotlin("kapt") version "1.9.10"

    application
}

group = "com.wittano"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

val picocliVersion = "4.7.4"
val jacksonVersion = "2.15.2"

dependencies {
    implementation("com.discord4j:discord4j-core:3.2.5")
    implementation("ch.qos.logback:logback-classic:1.4.11")
    implementation("org.codehaus.janino:janino:3.1.10")
    
    implementation("info.picocli:picocli:$picocliVersion")
    implementation("io.github.cdimascio:dotenv-kotlin:6.4.1")
    implementation("com.google.inject:guice:7.0.0")
    implementation("com.squareup.okhttp3:okhttp:4.11.0")
    implementation("io.projectreactor.kotlin:reactor-kotlin-extensions:1.2.2")

    implementation("com.fasterxml.jackson.core:jackson-core:$jacksonVersion")
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin:$jacksonVersion")

    kapt("info.picocli:picocli-codegen:$picocliVersion")

    testImplementation(kotlin("test"))
}

kapt {
    arguments {
        arg("project", "${project.group}/${project.name}")
    }
}

tasks.withType<Jar> {
    manifest {
        attributes["Main-Class"] = "com.wittano.komputer.MainKt"
    }

    duplicatesStrategy = DuplicatesStrategy.EXCLUDE

    from(sourceSets.main.get().output)

    dependsOn(configurations.runtimeClasspath)
    from({
        configurations.runtimeClasspath.get().filter { it.name.endsWith("jar") }.map { zipTree(it) }
    })
}

tasks.test {
    useJUnitPlatform()
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}

application {
    mainClass.set("com.wittano.komputer.MainKt")
}